import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { StorageClassCreateNfsComponent } from './storage-class-create-nfs.component';

describe('StorageClassCreateNfsComponent', () => {
  let component: StorageClassCreateNfsComponent;
  let fixture: ComponentFixture<StorageClassCreateNfsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ StorageClassCreateNfsComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(StorageClassCreateNfsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
