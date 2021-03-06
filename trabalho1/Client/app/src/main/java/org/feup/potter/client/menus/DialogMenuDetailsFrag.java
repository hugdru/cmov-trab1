package org.feup.potter.client.menus;

import android.app.Dialog;
import android.app.DialogFragment;
import android.os.Bundle;

import org.feup.potter.client.db.ItemInList;

public class DialogMenuDetailsFrag extends DialogFragment {
    public static DialogMenuDetailsFrag newInstance(ItemInList data) {
        DialogMenuDetailsFrag frag = new DialogMenuDetailsFrag();
        Bundle args = new Bundle();
        args.putSerializable("menuItem", data);
        frag.setArguments(args);
        return frag;
    }

    @Override
    public Dialog onCreateDialog(Bundle savedInstanceState) {
        ItemInList data = (ItemInList) getArguments().getSerializable("menuItem");

        DialogMenuDetails dialog = new DialogMenuDetails(getActivity(), data);
        return dialog;
    }
}
